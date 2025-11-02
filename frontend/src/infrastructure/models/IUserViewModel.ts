import type { IGuid } from '@/shared/Guid.ts';
import type { TUserRole } from '@/infrastructure/models/TUserRole.ts';
import type { JWTPayload } from 'jose';

export interface IUserViewModel extends JWTPayload {
  id: IGuid
  role: TUserRole
  token: string
  name: string
  surname: string
}
