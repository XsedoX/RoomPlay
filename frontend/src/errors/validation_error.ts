export default class ValidationError extends Error {
  public readonly fieldErrors: Record<string, string>;
  constructor(message: string, fieldErrors: Record<string, string>) {
    super(message);
    this.name = 'ValidationError';
    this.fieldErrors = fieldErrors;
  }
}
