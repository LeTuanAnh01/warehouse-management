import { AuthService } from './auth.service';
import { LoginDto, RegisterDto } from './auth.dto';
export declare class AuthController {
    private authService;
    constructor(authService: AuthService);
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
}
