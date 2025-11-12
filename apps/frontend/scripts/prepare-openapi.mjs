import { readFile, writeFile, mkdir, rm } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const responseRenameMap = {
  '400': 'Http400',
  '403': 'Http403',
  '404': 'Http404',
  '409': 'Http409',
  '500': 'Http500',
};

const schemaRenameOverrides = {
  AuthenticatedUserDoc: 'AuthenticatedUser',
  BackendErrorDoc: 'BackendError',
};

const sourceSpecPath = path.resolve(__dirname, '../../backend/ent/api/openapi.json');
const outputDir = path.resolve(__dirname, '../.orval');
const outputSpecPath = path.join(outputDir, 'openapi.json');
const generatedDir = path.resolve(__dirname, '../src/api/generated-orval');

const refRenameMap = {};

function renameResponses(responses) {
  if (!responses) return responses;
  return Object.entries(responses).reduce((acc, [key, value]) => {
    const nextKey = responseRenameMap[key] ?? key;
    acc[nextKey] = value;
    if (nextKey !== key) {
      refRenameMap[`#/components/responses/${key}`] = `#/components/responses/${nextKey}`;
    }
    return acc;
  }, {});
}

function renameSchemas(schemas) {
  if (!schemas) return schemas;

  const result = {};
  for (const [key, value] of Object.entries(schemas)) {
    const base = key.includes('.') ? key.slice(key.lastIndexOf('.') + 1) : key;
    const nextKey = schemaRenameOverrides[base] ?? base;

    if (result[nextKey]) {
      // Preserve the existing entry if a duplicate sanitized name appears.
      continue;
    }

    result[nextKey] = value;
    if (nextKey !== key) {
      refRenameMap[`#/components/schemas/${key}`] = `#/components/schemas/${nextKey}`;
    }
  }

  return result;
}

function updateRefs(node) {
  if (Array.isArray(node)) {
    node.forEach(updateRefs);
    return;
  }

  if (node && typeof node === 'object') {
    for (const [key, value] of Object.entries(node)) {
      if (key === '$ref' && typeof value === 'string') {
        node[key] = refRenameMap[value] ?? value;
      } else {
        updateRefs(value);
      }
    }
  }
}

try {
  const raw = await readFile(sourceSpecPath, 'utf-8');
  const spec = JSON.parse(raw);

  spec.components = spec.components ?? {};
  spec.components.responses = renameResponses(spec.components.responses);
  spec.components.schemas = renameSchemas(spec.components.schemas);

  updateRefs(spec);

  await mkdir(outputDir, { recursive: true });
  await rm(generatedDir, { recursive: true, force: true });
  await writeFile(outputSpecPath, JSON.stringify(spec, null, 2), 'utf-8');
  console.log(`Prepared OpenAPI spec with sanitized response names at ${outputSpecPath}`);
} catch (error) {
  console.error('Failed to prepare OpenAPI spec for Orval:', error);
  process.exit(1);
}
