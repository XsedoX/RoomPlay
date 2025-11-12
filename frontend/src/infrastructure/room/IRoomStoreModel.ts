import type { TUserRole } from '@/infrastructure/user/TUserRole.ts'

export default interface IRoomStoreModel {
  name: string;
  boostUsedAtUtc: Date | null;
  qrCode: string;
  userRole: TUserRole
}
