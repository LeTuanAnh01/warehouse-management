import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from './user.entity';

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
  ) {}

  findAll() {
    return this.usersRepository.find({ select: ['id', 'username', 'email', 'fullName', 'phone', 'isActive', 'roleId', 'lastLogin', 'createdAt', 'updatedAt'] });
  }

  async findOne(id: number) {
    const user = await this.usersRepository.findOne({ where: { id }, select: ['id', 'username', 'email', 'fullName', 'phone', 'isActive', 'roleId', 'lastLogin', 'createdAt', 'updatedAt'] });
    if (!user) throw new NotFoundException('User not found');
    return user;
  }

  async update(id: number, dto: Partial<Pick<User, 'fullName' | 'phone' | 'isActive'>>) {
    await this.usersRepository.update(id, dto);
    return this.findOne(id);
  }

  async remove(id: number) {
    const user = await this.findOne(id);
    await this.usersRepository.remove(user);
    return { message: 'User deleted' };
  }
}
