import { Role } from './role.entity';
export declare class User {
    id: number;
    username: string;
    email: string;
    passwordHash: string;
    fullName: string;
    phone: string;
    roleId: number;
    isActive: boolean;
    lastLogin: Date;
    createdAt: Date;
    updatedAt: Date;
    role: Role;
}
