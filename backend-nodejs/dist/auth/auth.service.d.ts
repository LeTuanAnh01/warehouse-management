import { Repository } from 'typeorm';
import { JwtService } from '@nestjs/jwt';
import { User } from '../users/user.entity';
import { LoginDto, RegisterDto } from './auth.dto';
export declare class AuthService {
    private usersRepository;
    private jwtService;
    constructor(usersRepository: Repository<User>, jwtService: JwtService);
    register(dto: RegisterDto): Promise<{
        id: number;
        username: string;
        email: string;
        fullName: string;
        phone: string;
        roleId: number;
        isActive: boolean;
        lastLogin: Date;
        createdAt: Date;
        updatedAt: Date;
        role: import("../users/role.entity").Role;
    }>;
    login(dto: LoginDto): Promise<{
        access_token: string;
        user: {
            id: number;
            username: string;
            email: string;
            fullName: string;
            phone: string;
            roleId: number;
            isActive: boolean;
            lastLogin: Date;
            createdAt: Date;
            updatedAt: Date;
            role: import("../users/role.entity").Role;
        };
    }>;
    validateUser(userId: number): Promise<User | null>;
}
