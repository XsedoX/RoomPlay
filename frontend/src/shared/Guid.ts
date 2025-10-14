export class Guid implements IGuid {
  private readonly value: string;
  constructor(value: string) {
    if (!Guid.isValid(value)) {
      throw new Error(`Invalid GUID format: ${value}`);
    }
    this.value = value;
  }
  static generate(): Guid {
    const guidString = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = Math.random() * 16 | 0;
      const v = c === 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
    return new Guid(guidString);
  }
  static isValid(value: string): boolean {
    const regex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    return regex.test(value);
  }
  toString(): string {
    return this.value;
  }
  valueOf(): string {
    return this.value;
  }
  equals(other: Guid): boolean {
    return this.value === other.value;
  }
}
export interface IGuid{
  toString(): string;
  equals(other: Guid): boolean;
}

