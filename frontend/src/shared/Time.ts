export interface ITime {
  hour: number,
  minute: number
  toString(): string
  toDate(): Date
}
export class Time implements ITime {
  hour: number;
  minute: number;

  constructor();
  constructor(date: Date);
  constructor(hour: number, minute: number)
  constructor(timeString: string);
  constructor(hourOrDateOrString?: number | Date | string, minute?: number) {
    if(hourOrDateOrString instanceof Date){
      this.hour = hourOrDateOrString.getHours();
      this.minute = hourOrDateOrString.getMinutes();
    }
    else if(typeof hourOrDateOrString === 'string'){
      const parts = hourOrDateOrString.split(':');
      if (parts.length !== 2) {
        throw new Error('Invalid time string format. Expected "HH:mm".');
      }
      const hour = Number(parts[0]);
      const min = Number(parts[1]);

      if (
        !Number.isNaN(hour) && hour >= 0 && hour <= 23 &&
        !Number.isNaN(min) && min >= 0 && min <= 59
      ) {
        this.hour = hour;
        this.minute = min;
      } else {
        throw new Error('Invalid time string format. Expected "HH:mm".');
      }
    }
    else if(typeof hourOrDateOrString === 'number' && typeof minute === 'number'){
      this.hour = hourOrDateOrString;
      this.minute = minute;
    }
    else{
      const now = new Date();
      this.hour = now.getHours();
      this.minute = now.getMinutes();
    }
  }

  toString(): string {
    const pad = (num: number) => num.toString().padStart(2, '0');
    return `${pad(this.hour)}:${pad(this.minute)}`;
  }
  toDate(): Date {
    const date = new Date();
    date.setHours(this.hour, this.minute, 0, 0);
    return date;
  }
}
