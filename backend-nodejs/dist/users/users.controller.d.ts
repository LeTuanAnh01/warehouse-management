import { UsersService } from './users.service';
export declare class UsersController {
    private usersService;
    constructor(usersService: UsersService);
    getProfile(req: any): Promise<import("./user.entity").User>;
    findAll(): Promise<import("./user.entity").User[]>;
    findOne(id: string): Promise<import("./user.entity").User>;
    update(id: string, dto: any): Promise<import("./user.entity").User>;
    remove(id: string): Promise<{
        message: string;
    }>;
}
