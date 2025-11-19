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

const sourceSpecPath = path.resolve(__dirname, '../../backend/api/swag-docs/swagger.json');
const outputDir = path.resolve(__dirname, '../.orval');
const outputSpecPath = path.join(outputDir, 'openapi.json');
const generatedDir = path.resolve(__dirname, '../src/api/generated-orval');

const refRenameMap = {};

function convertSwagger2ToOpenAPI3(swagger) {
  const openapi = {
    openapi: '3.0.3',
    info: swagger.info,
    servers: [
      {
        url: `http://${swagger.host}${swagger.basePath || ''}`,
        description: 'Development server',
      },
    ],
    paths: {},
    components: {
      schemas: swagger.definitions || {},
      securitySchemes: {},
    },
  };

  // Convert security definitions
  if (swagger.securityDefinitions) {
    for (const [name, def] of Object.entries(swagger.securityDefinitions)) {
      if (def.type === 'apiKey') {
        openapi.components.securitySchemes[name] = {
          type: 'apiKey',
          in: def.in,
          name: def.name,
          description: def.description,
        };
      }
    }
  }

  // Convert paths
  if (swagger.paths) {
    for (const [path, pathItem] of Object.entries(swagger.paths)) {
      openapi.paths[path] = {};
      for (const [method, operation] of Object.entries(pathItem)) {
        if (['get', 'post', 'put', 'delete', 'patch', 'options', 'head'].includes(method)) {
          const newOperation = {
            summary: operation.summary,
            description: operation.description,
            operationId: operation.operationId,
            tags: operation.tags,
            security: operation.security,
            responses: {},
          };

          // Convert parameters
          if (operation.parameters) {
            const pathParams = [];
            const queryParams = [];
            const headerParams = [];
            let bodyParam = null;

            for (const param of operation.parameters) {
              if (param.in === 'body') {
                bodyParam = param;
              } else if (param.in === 'path') {
                pathParams.push({
                  name: param.name,
                  in: 'path',
                  required: param.required !== false,
                  description: param.description,
                  schema: param.type ? { type: param.type, format: param.format } : param.schema || { type: 'string' },
                });
              } else if (param.in === 'query') {
                queryParams.push({
                  name: param.name,
                  in: 'query',
                  required: param.required === true,
                  description: param.description,
                  schema: param.type ? { type: param.type, format: param.format } : param.schema || { type: 'string' },
                });
              } else if (param.in === 'header') {
                headerParams.push({
                  name: param.name,
                  in: 'header',
                  required: param.required === true,
                  description: param.description,
                  schema: param.type ? { type: param.type } : param.schema || { type: 'string' },
                });
              }
            }

            const allParams = [...pathParams, ...queryParams, ...headerParams];
            if (allParams.length > 0) {
              newOperation.parameters = allParams;
            }

            if (bodyParam) {
              newOperation.requestBody = {
                required: bodyParam.required !== false,
                content: {
                  'application/json': {
                    schema: bodyParam.schema,
                  },
                },
              };
            }
          }

          // Convert responses
          for (const [code, response] of Object.entries(operation.responses || {})) {
            newOperation.responses[code] = {
              description: response.description || code === '200' ? 'Success' : 'Error',
            };
            if (response.schema) {
              newOperation.responses[code].content = {
                'application/json': {
                  schema: response.schema,
                },
              };
            }
          }

          openapi.paths[path][method] = newOperation;
        }
      }
    }
  }

  // Update all $ref paths from #/definitions/ to #/components/schemas/
  const updateRefs = (obj) => {
    if (Array.isArray(obj)) {
      obj.forEach(updateRefs);
    } else if (obj && typeof obj === 'object') {
      for (const [key, value] of Object.entries(obj)) {
        if (key === '$ref' && typeof value === 'string') {
          obj[key] = value.replace('#/definitions/', '#/components/schemas/');
        } else {
          updateRefs(value);
        }
      }
    }
  };

  updateRefs(openapi);
  return openapi;
}

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
    // Take the last token after the last dot and capitalize it
    const base = key.includes('.') ? key.slice(key.lastIndexOf('.') + 1) : key;
    const nextKey = base.charAt(0).toUpperCase() + base.slice(1);

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
  let spec = JSON.parse(raw);

  // Convert Swagger 2.0 to OpenAPI 3.0 if needed
  if (spec.swagger === '2.0') {
    spec = convertSwagger2ToOpenAPI3(spec);
  }

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
