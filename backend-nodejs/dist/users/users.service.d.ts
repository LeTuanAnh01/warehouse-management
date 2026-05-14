import { Repository } from 'typeorm';
import { User } from './user.entity';
export declare class UsersService {
    private usersRepository;
    constructor(usersRepository: Repository<User>);
    findAll(): Promise<User[]>;
    findOne(id: number): Promise<User>;
    update(id: number, dto: Partial<Pick<User, 'fullName' | 'phone' | 'isActive'>>): Promise<User>;
    remove(id: number): Promise<{
        message: string;
    }>;
}
