export class Doi {
  private readonly _doi: string;

  private constructor(doi: string) {
    this._doi = doi;
  }

  // Matches standard DOI formats: 10.xxxx/yyyy
  private static readonly ModernDoiRegex =
    /^10\.\d{4,9}\/[-._;()/:a-z0-9<>]+$/i;
  // Matches older Wiley DOI formats: 10.1002/xxxx
  private static readonly OldDoiRegex = /^10\.1002\/[^\s]+$/i;

  /**
   * Try to parse a DOI string into a Doi object.
   * @param doi The DOI string to parse.
   * @returns The Doi object if successful, otherwise null.
   */
  public static tryParse(doi: string): Doi | null {
    try {
      const normalized = Doi.normalize(doi);
      if (Doi.isValid(normalized)) {
        return new Doi(normalized);
      }
    } catch {
      // Ignore errors during normalization
    }
    return null;
  }

  /**
   * Parse a DOI string into a Doi object.
   * @param doi The DOI string to parse.
   * @returns The Doi object.
   * @throws Error if the DOI is invalid.
   */
  public static parse(doi: string): Doi {
    const result = Doi.tryParse(doi);
    if (result) {
      return result;
    }
    throw new Error(`Invalid DOI: ${doi}`);
  }

  /**
   * Normalize a DOI string.
   * @param doi The DOI string to normalize.
   * @returns The normalized DOI string.
   * @throws Error if the URL is invalid or not a doi.org link.
   */
  public static normalize(doi: string): string {
    doi = doi.trim();
    if (!doi.startsWith("http://") && !doi.startsWith("https://")) {
      return doi.replace(/^\/+|\/+$/g, ""); // Trim slashes
    }

    try {
      const url = new URL(doi);
      if (url.hostname !== "doi.org" && url.hostname !== "dx.doi.org") {
        throw new Error("Missing doi.org from URL");
      }
      return url.pathname.replace(/^\/+|\/+$/g, "");
    } catch (e) {
      throw new Error("Invalid URL");
    }
  }

  /**
   * Check if a DOI string is valid.
   * @param doi The DOI string to check.
   * @returns True if the DOI is valid, otherwise false.
   */
  public static isValid(doi: string): boolean {
    return Doi.ModernDoiRegex.test(doi) || Doi.OldDoiRegex.test(doi);
  }

  public toString(): string {
    return this._doi;
  }

  public toUrl(): string {
    return `https://doi.org/${this._doi}`;
  }
}
