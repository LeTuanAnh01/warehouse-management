import { Injectable, UnauthorizedException, ConflictException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
import { User } from '../users/user.entity';
import { LoginDto, RegisterDto } from './auth.dto';

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
    private jwtService: JwtService,
  ) {}

  async register(dto: RegisterDto) {
    const exists = await this.usersRepository.findOne({
      where: [{ username: dto.username }, { email: dto.email }],
    });
    if (exists) {
      throw new ConflictException('Username or email already exists');
    }

    const passwordHash = await bcrypt.hash(dto.password, 10);
    const user = this.usersRepository.create({
      username: dto.username,
      email: dto.email,
      passwordHash,
      fullName: dto.fullName,
      phone: dto.phone,
      roleId: 2,
    });
    await this.usersRepository.save(user);

    const { passwordHash: _, ...result } = user;
    return result;
  }

  async login(dto: LoginDto) {
    const user = await this.usersRepository.findOne({
      where: { username: dto.username },
    });
    if (!user || !(await bcrypt.compare(dto.password, user.passwordHash))) {
      throw new UnauthorizedException('Invalid credentials');
    }
    if (!user.isActive) {
      throw new UnauthorizedException('Account is disabled');
    }

    await this.usersRepository.update(user.id, { lastLogin: new Date() });

    const payload = { sub: user.id, username: user.username, roleId: user.roleId };
    const token = this.jwtService.sign(payload);

    const { passwordHash: _, ...userInfo } = user;
    return { access_token: token, user: userInfo };
  }

  async validateUser(userId: number): Promise<User | null> {
    return this.usersRepository.findOne({ where: { id: userId } });
  }
}
