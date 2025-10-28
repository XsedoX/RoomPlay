import type { IGuid } from '@/shared/Guid.ts';
import type { TUserRole } from '@/repositories/TUserRole.ts';
import type { JWTPayload } from 'jose';

export interface IUserViewModel extends JWTPayload {
  id: IGuid
  role: TUserRole
  token: string
  name: string
  surname: string
}
