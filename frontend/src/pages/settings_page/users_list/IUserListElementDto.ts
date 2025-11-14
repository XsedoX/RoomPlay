import type { IGuid } from '@/shared/Guid.ts';

export default interface IUserListElementDto {
  name: string;
  surname: string;
  id: IGuid
}
