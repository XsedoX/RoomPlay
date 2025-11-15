import type { TUserRole } from '@/infrastructure/user/TUserRole.ts';

export default interface IRoomStoreModel {
  name: string;
  boostData: IBoostData | null;
  qrCode: string;
  userRole: TUserRole;
}
export interface IBoostData {
  boostUsedAtUtc: Date;
  boostCooldownSeconds: number;
}
