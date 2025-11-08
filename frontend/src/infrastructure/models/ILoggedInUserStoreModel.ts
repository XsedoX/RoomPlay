import { TUserRole } from '@/infrastructure/models/TUserRole.ts';
import type { IGuid } from '@/shared/Guid.ts';

export interface ILoggedInUserStoreModel {
  name: string;
  surname: string;
  roomId: IGuid | null;
  role: TUserRole | null;
}
