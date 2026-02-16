/**
 * EU VAT number parsing and validation.
 *
 * VAT numbers consist of a 2-letter country code followed by a number.
 * Each country has its own format, but this provides basic validation.
 *
 * Examples:
 * - NL822655287B01 (Netherlands)
 * - DE123456789 (Germany)
 * - BE0123456789 (Belgium)
 */
export class Vat {
  private readonly _countryCode: string;
  private readonly _number: string;
  private readonly _raw: string;

  private constructor(countryCode: string, number: string, raw: string) {
    this._countryCode = countryCode;
    this._number = number;
    this._raw = raw;
  }

  // EU country codes that participate in VIES
  private static readonly EuCountryCodes = new Set([
    "AT", // Austria
    "BE", // Belgium
    "BG", // Bulgaria
    "CY", // Cyprus
    "CZ", // Czech Republic
    "DE", // Germany
    "DK", // Denmark
    "EE", // Estonia
    "EL", // Greece (uses EL instead of GR)
    "ES", // Spain
    "FI", // Finland
    "FR", // France
    "HR", // Croatia
    "HU", // Hungary
    "IE", // Ireland
    "IT", // Italy
    "LT", // Lithuania
    "LU", // Luxembourg
    "LV", // Latvia
    "MT", // Malta
    "NL", // Netherlands
    "PL", // Poland
    "PT", // Portugal
    "RO", // Romania
    "SE", // Sweden
    "SI", // Slovenia
    "SK", // Slovakia
    "XI", // Northern Ireland (post-Brexit)
  ]);

  // Basic VAT format: 2-letter country code + alphanumeric number
  private static readonly VatRegex = /^([A-Z]{2})([A-Z0-9]{2,14})$/i;

  /**
   * Try to parse a VAT number string into a Vat object.
   * @param vat The VAT number string to parse.
   * @returns The Vat object if successful, otherwise null.
   */
  public static tryParse(vat: string): Vat | null {
    try {
      const normalized = Vat.normalize(vat);
      const match = normalized.match(Vat.VatRegex);

      if (match) {
        const countryCode = match[1].toUpperCase();
        const number = match[2].toUpperCase();

        if (Vat.EuCountryCodes.has(countryCode)) {
          return new Vat(countryCode, number, normalized);
        }
      }
    } catch {
      // Ignore errors during parsing
    }
    return null;
  }

  /**
   * Parse a VAT number string into a Vat object.
   * @param vat The VAT number string to parse.
   * @returns The Vat object.
   * @throws Error if the VAT is invalid.
   */
  public static parse(vat: string): Vat {
    const result = Vat.tryParse(vat);
    if (result) {
      return result;
    }
    throw new Error(`Invalid VAT number: ${vat}`);
  }

  /**
   * Normalize a VAT number string.
   * Removes spaces, dots, dashes and converts to uppercase.
   * @param vat The VAT number string to normalize.
   * @returns The normalized VAT number string.
   */
  public static normalize(vat: string): string {
    return vat
      .trim()
      .toUpperCase()
      .replace(/[\s.\-]/g, "");
  }

  /**
   * Check if a VAT number string appears to be valid.
   * This only checks format, not validity via VIES API.
   * @param vat The VAT number string to check.
   * @returns True if the VAT appears valid, otherwise false.
   */
  public static isValid(vat: string): boolean {
    return Vat.tryParse(vat) !== null;
  }

  /** The 2-letter country code. */
  public get countryCode(): string {
    return this._countryCode;
  }

  /** The VAT number without country code. */
  public get number(): string {
    return this._number;
  }

  /** Full VAT number including country code. */
  public toString(): string {
    return this._raw;
  }

  /** URL to VIES VAT validation page. */
  public toUrl(): string {
    return `https://ec.europa.eu/taxation_customs/vies/#/vat-validation`;
  }
}
