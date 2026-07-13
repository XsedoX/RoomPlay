import type { IGuid } from '@/shared/guid/Guid';

export default interface IUserListElementDto {
  name: string;
  surname: string;
  id: IGuid;
}
