export interface ITime {
  hour: number;
  minute: number;
  seconds: number;
  toString(): string;
  toDate(): Date;
}
export class Time implements ITime {
  hour: number;
  minute: number;
  seconds = 0;

  constructor();
  constructor(date: Date);
  constructor(hour: number, minute: number, seconds?: number);
  constructor(timeString: string);
  constructor(hourOrDateOrString?: number | Date | string, minute?: number, seconds?: number) {
    if (hourOrDateOrString instanceof Date) {
      this.hour = hourOrDateOrString.getHours();
      this.minute = hourOrDateOrString.getMinutes();
      this.seconds = hourOrDateOrString.getSeconds();
    } else if (typeof hourOrDateOrString === 'string') {
      const parts = hourOrDateOrString.split(':');
      if (parts.length !== 2 && parts.length !== 3) {
        throw new Error('Invalid time string format. Expected "HH:mm:ss" or "HH:mm".');
      }
      const hour = Number(parts[0]);
      const min = Number(parts[1]);
      let secs = 0;
      if (parts.length === 3) {
        secs = Number(parts[2]);
      }
      if (
        !Number.isNaN(hour) &&
        hour >= 0 &&
        hour <= 23 &&
        !Number.isNaN(min) &&
        min >= 0 &&
        min <= 59
      ) {
        this.hour = hour;
        this.minute = min;
        if (!Number.isNaN(secs) && secs >= 0 && secs <= 59) {
          this.seconds = secs;
        }
      } else {
        throw new Error('Invalid time string format. Expected "HH:mm:ss" or "HH:mm".');
      }
    } else if (typeof hourOrDateOrString === 'number' && typeof minute === 'number') {
      this.hour = hourOrDateOrString;
      this.minute = minute;
      if (typeof seconds === 'number') {
        this.seconds = seconds;
      }
    } else {
      const now = new Date();
      this.hour = now.getHours();
      this.minute = now.getMinutes();
      this.seconds = now.getSeconds();
    }
  }
  totalSeconds(): number {
    return this.hour * 3600 + this.minute * 60 + this.seconds;
  }
  static fromSeconds(totalSeconds: number): Time {
    const hours = Math.floor(totalSeconds / 3600) % 24;
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    return new Time(hours, minutes, seconds);
  }
  static from(startTime: Date | string = new Date()): TimeDurationBuilder {
    return new TimeDurationBuilder(startTime);
  }
  toString(): string {
    const pad = (num: number) => num.toString().padStart(2, '0');
    if (this.hour === 0 && this.minute === 0) return this.seconds.toString();
    if (this.hour === 0) return `${this.minute}:${pad(this.seconds)}`;
    return `${pad(this.hour)}:${pad(this.minute)}:${pad(this.seconds)}`;
  }
  toDate(): Date {
    const date = new Date();
    date.setHours(this.hour, this.minute, this.seconds, 0);
    return date;
  }
  incrementSeconds(seconds: number = 1): Time {
    const newSeconds = this.totalSeconds() + seconds;
    const newTime = Time.fromSeconds(newSeconds);
    this.hour = newTime.hour;
    this.minute = newTime.minute;
    this.seconds = newTime.seconds;
    return this;
  }
}
class TimeDurationBuilder {
  private readonly startDate: Date;

  constructor(startTime: Date | string) {
    this.startDate = startTime instanceof Date ? startTime : new Date(startTime);
  }

  to(endTime: Date | string = new Date()): Time {
    const endDate = endTime instanceof Date ? endTime : new Date(endTime);
    const differenceInMs = endDate.getTime() - this.startDate.getTime();
    const differenceInSeconds = Math.floor(differenceInMs / 1000);
    return Time.fromSeconds(differenceInSeconds);
  }
}
