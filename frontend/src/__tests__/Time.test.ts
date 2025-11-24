import { describe, it, expect, vi } from 'vitest';
import { Time } from '../shared/Time';

describe('Time Class', () => {
  describe('Constructors', () => {
    it('should create a Time instance with current time by default', () => {
      const date = new Date(2023, 9, 10, 14, 30, 45);
      vi.useFakeTimers();
      vi.setSystemTime(date);

      const time = new Time();
      expect(time.hour).toBe(14);
      expect(time.minute).toBe(30);
      expect(time.seconds).toBe(45);

      vi.useRealTimers();
    });

    it('should create a Time instance from a Date object', () => {
      const date = new Date(2023, 9, 10, 10, 15, 30);
      const time = new Time(date);
      expect(time.hour).toBe(10);
      expect(time.minute).toBe(15);
      expect(time.seconds).toBe(30);
    });

    it('should create a Time instance from hour, minute, and seconds', () => {
      const time = new Time(12, 45, 10);
      expect(time.hour).toBe(12);
      expect(time.minute).toBe(45);
      expect(time.seconds).toBe(10);
    });

    it('should create a Time instance from hour and minute (seconds default to 0)', () => {
      const time = new Time(12, 45);
      expect(time.hour).toBe(12);
      expect(time.minute).toBe(45);
      expect(time.seconds).toBe(0); // Assuming default is handled or undefined if optional but class init says seconds=0
    });

    it('should create a Time instance from string "HH:mm:ss"', () => {
      const time = new Time('14:30:15');
      expect(time.hour).toBe(14);
      expect(time.minute).toBe(30);
      expect(time.seconds).toBe(15);
    });

    it('should create a Time instance from string "HH:mm"', () => {
      const time = new Time('14:30');
      expect(time.hour).toBe(14);
      expect(time.minute).toBe(30);
      expect(time.seconds).toBe(0);
    });

    it('should throw error for invalid string "invalid"', () => {
      expect(() => new Time('invalid')).toThrow(
        'Invalid time string format. Expected "HH:mm:ss" or "HH:mm".',
      );
    });

    it('should throw error for invalid hour "25:00:00"', () => {
      expect(() => new Time('25:00:00')).toThrow(
        'Invalid time string format. Expected "HH:mm:ss" or "HH:mm".',
      );
    });

    it('should throw error for invalid minute "12:60:00"', () => {
      expect(() => new Time('12:60:00')).toThrow(
        'Invalid time string format. Expected "HH:mm:ss" or "HH:mm".',
      );
    });

    it('should throw error for invalid seconds "12:00:60"', () => {
      expect(() => new Time('12:00:60')).toThrow(
        'Invalid time string format. Expected "HH:mm:ss" or "HH:mm".',
      );
    });
  });

  describe('Methods', () => {
    it('totalSeconds should return correct total seconds', () => {
      const time = new Time(1, 1, 1);
      expect(time.totalSeconds()).toBe(3600 + 60 + 1);
    });

    it('toString should return formatted string "HH:mm:ss"', () => {
      const time = new Time(9, 5, 3);
      expect(time.toString()).toBe('09:05:03');
    });

    it('toString should return formatted string "mm:ss" if hour is 0', () => {
      const time = new Time(0, 45, 30);
      expect(time.toString()).toBe('45:30');
    });

    it('toString should return seconds if hour and minute are 0', () => {
      const time = new Time(0, 0, 45);
      expect(time.toString()).toBe('45');
    });

    it('toDate should return a Date object with correct time', () => {
      const time = new Time(14, 30, 0);
      const date = time.toDate();
      expect(date.getHours()).toBe(14);
      expect(date.getMinutes()).toBe(30);
      expect(date.getSeconds()).toBe(0);
    });

    it('incrementSeconds should increment time correctly', () => {
      const time = new Time(10, 0, 0);
      time.incrementSeconds(65);
      expect(time.hour).toBe(10);
      expect(time.minute).toBe(1);
      expect(time.seconds).toBe(5);
    });

    it('incrementSeconds should handle day overflow (wrap around)', () => {
      // The implementation of fromSeconds uses % 24 for hours, so it wraps around.
      const time = new Time(23, 59, 59);
      time.incrementSeconds(1);
      expect(time.hour).toBe(0);
      expect(time.minute).toBe(0);
      expect(time.seconds).toBe(0);
    });
  });

  describe('Static Methods', () => {
    it('fromSeconds should create Time instance from total seconds', () => {
      const totalSeconds = 3665; // 1h 1m 5s
      const time = Time.fromSeconds(totalSeconds);
      expect(time.hour).toBe(1);
      expect(time.minute).toBe(1);
      expect(time.seconds).toBe(5);
    });

    it('from().to() should calculate duration between two times', () => {
      const start = new Date(2023, 0, 1, 10, 0, 0);
      const end = new Date(2023, 0, 1, 12, 30, 15);

      const duration = Time.from(start).to(end);
      expect(duration.hour).toBe(2);
      expect(duration.minute).toBe(30);
      expect(duration.seconds).toBe(15);
    });

    it('from().to() should use current time if end time is not provided', () => {
      const start = new Date(2023, 0, 1, 10, 0, 0);
      const now = new Date(2023, 0, 1, 11, 0, 0);

      vi.useFakeTimers();
      vi.setSystemTime(now);

      const duration = Time.from(start).to();
      expect(duration.hour).toBe(1);
      expect(duration.minute).toBe(0);
      expect(duration.seconds).toBe(0);

      vi.useRealTimers();
    });
  });
});
