!  2�����L���r�e�B�t���[�̐��l�v�Z
!  2022.07.20
!  ���͂̃|�A�\���������̒��̍����𒆐S�������畗�㍷����

        implicit real*8(a-h, o-z), integer*4(i-n)
        character file_name*7
        parameter (mm=260, nn=260)
        dimension u(0:mm,0:nn), unew(0:mm,0:nn)
        dimension v(0:mm,0:nn), vnew(0:mm,0:nn)
        dimension p(0:mm,0:nn), pnew(0:mm,0:nn)
        dimension u4(0:mm,0:nn), v4(0:mm,0:nn)
        dimension uat1(0:mm,0:nn), uat2(0:mm,0:nn)
        dimension upt(0:mm,0:nn),udt(0:mm,0:nn)
        dimension vat1(0:mm,0:nn), vat2(0:mm,0:nn)
        dimension vpt(0:mm,0:nn),vdt(0:mm,0:nn)
        dimension pp1(0:mm,0:nn), pp2(0:mm,0:nn)
        dimension pp3(0:mm,0:nn),phi(0:mm,0:nn)
        dimension uout(0:mm,0:nn), vout(0:mm,0:nn)
		
        open(10,file = 'fname.prn')  ! �J�^���O�t�@�C��
        
!       �v�Z����
        mx = 64
        my = 64
        xL = 0.02d0             ! �L���r�e�B�T�C�Y [2cm]
        ktend = 10000            ! ���Ԕ��W�v�Z�̉�
        dh = xL / dble(mx)      ! ���b�V���T�C�Y

        dt = 0.001d0			! [sec]
        Us = 0.02d0  			! [m/s] �L���r�e�B�̊O�̎嗬�̗���
        eps = 10.d0**(-6)       ! ���͂̎����v�Z��臒l
        kend = 100000	 	    ! ���͂̎����v�Z�̌��E��
        dvs = 10.d0**(-6)       ! ���̓��S���W�� [m^2/s]
        rho = 1000.d0           ! ���̖��x [kg/m^3]
		alpha = 0.5d0
 
 !      ��������
        do 10 i = 1, mx+1	! ���������̗���̌v�Z�_�́C1�`mx+1
            do 10 j = 1, my
                u(i,j) = 0.d0
   10   continue
     
        do 11 i = 1, mx
            do 11 j = 1, my+1	! ���f�����̗���̌v�Z�_�́C1�`my+1
                v(i,j) = 0.d0
   11   continue
     
        do 12 i = 1, mx
            do 12 j = 1, my		! ���f�����̈��͂̌v�Z�_�D�嗬�̒��͌v�Z���Ȃ�
                p(i,j) = 0.4d0	! Re = 400�̎��̏����l
   12   continue
   
! ----- ���Ԕ��W�̌v�Z�i���C���j------------------------------------
        do 500 kt = 1, ktend
!        write(*,*)
!        write(*,*) 'kt = ',kt

!       ���E����
!			�嗬�̒��̈��͂͌v�Z���Ȃ��D
!			�嗬�̉����̉������x���v�Z���Ȃ��D

        do 20 i = 1, mx
            u(i,my+1) = Us			! �嗬�̒��̒l
            v(i,my+1) = 0.d0		! �嗬�ւ͗��ꍞ�܂Ȃ�
            v(i,my+2) = 0.d0		! �嗬�̏㑤�̕ǂ̒l
            p(i,my+1) = p(i,my)		! �嗬�̒��̈��͂͌v�Z���Ȃ�
            p(i,my+2) = p(i,my+1)	! �嗬�̂���ɏ㑤�̒l��^����
            u(i,0) = -u(i,1)			! �ǂ̉���
            v(i,1) = 0.d0			! �����̕ǂ̒l
            p(i,0) = p(i,1)			! �ǂ̉���
   20   continue

        do 21 j = 1, my
            u(1,j) = 0.d0			! ��
            v(0,j) = -v(1,j)
            p(0,j) = p(1,j)
            u(mx+1,j) = 0.d0		! �E
            v(mx+1,j) = -v(mx,j)
            p(mx+1,j) = p(mx,j)
   21   continue
   
!       u�̒�`�_��v4�Cv�̒�`�_�ɂ�����u4
        do 22 i = 1, mx
            do 22 j = 2, my
                u4(i,j) = (u(i,j) + u(i+1,j) + u(i,j-1) + u(i+1,j-1)) * 0.25d0
   22   continue
   
        do 23 i = 2, mx
            do 23 j = 1, my
                v4(i,j) = (v(i,j) + v(i-1,j) + v(i,j+1) + v(i-1,j+1)) * 0.25d0
   23   continue
           
!       NS�������̌v�Z(unew)
        do 24 i = 2, mx
            do 24 j = 1, my
                uat1(i,j) =  ( u(i,j) + dabs( u(i,j)))/2.d0 * (u(i,j)   - u(i-1,j))/dh &
                          +  ( u(i,j) - dabs( u(i,j)))/2.d0 * (u(i+1,j) - u(i,j))/dh

                uat2(i,j) =  (v4(i,j) + dabs(v4(i,j)))/2.d0 * (u(i,j)   - u(i,j-1))/dh &
                          +  (v4(i,j) - dabs(v4(i,j)))/2.d0 * (u(i,j+1) - u(i,j))/dh

                upt(i,j)  = (p(i,j) - p(i-1,j))/(dh * rho)

                udt(i,j)  = (u(i+1,j) + u(i-1,j) + u(i,j-1) + u(i,j+1) - 4.d0 * u(i,j)) / (dh**2) * dvs

                unew(i,j) = u(i,j) - dt*(uat1(i,j) + uat2(i,j) +upt(i,j) - udt(i,j))
   24   continue

!       NS�������̌v�Z(vnew)
        do 25 i = 1, mx
            do 25 j = 2, my        
                vat1(i,j) =  (u4(i,j) + dabs(u4(i,j)))/2.d0 * (v(i,j)   - v(i-1,j))/dh &
                          +  (u4(i,j) - dabs(u4(i,j)))/2.d0 * (v(i+1,j) - v(i,j))/dh
                  
                vat2(i,j) =  ( v(i,j) + dabs( v(i,j)))/2.d0 * (v(i,j)   - v(i,j-1))/dh &
                          +  ( v(i,j) - dabs( v(i,j)))/2.d0 * (v(i,j+1) - v(i,j))/dh
                  
                vpt(i,j)  = (p(i,j) - p(i,j-1))/(dh * rho)
 
                vdt(i,j)  = (v(i+1,j) + v(i-1,j) + v(i,j-1) + v(i,j+1) - 4.d0 * v(i,j)) / (dh**2) * dvs
        
                vnew(i,j) = v(i,j) - dt*(vat1(i,j) + vat2(i,j) +vpt(i,j) - vdt(i,j))
   25   continue
   
        do 26 i = 2, mx
            do 26 j = 1, my
                u(i,j) = unew(i,j)   ! NS�������ɂ���ĐV�����Ȃ���unew��u�ɏ㏑������ �� u�̎�������i��
   26   continue
   
        do 27 i = 1, mx
            do 27 j = 2, my
                v(i,j) = vnew(i,j)   ! NS�������ɂ���ĐV�����Ȃ���vnew��v�ɏ㏑������ �� v�̎�������i��
   27   continue
   
!------------------------------------------------------------------- �����܂ł�u��v���X�V���ꂽ

!       ���E�����i�L���r�e�B�̒����ꎞ���i�񂾂̂ŁC���E���V��������j
        do 30 i = 1, mx
            u(i,my+1) = Us			! ��
            u(i,my+2) = u(i,my+1)
            v(i,my+1) = 0.d0
            v(i,my+2) = v(i,my+1)
            v(i,my+3) = v(i,my+2)
            p(i,my+1) = p(i,my)
            p(i,my+2) = p(i,my+1)
            u(i,0) = -u(i,1)			! ��
            v(i,1) = 0.d0
            v(i,0) = v(i,1)
            p(i,0) = p(i,1)
   30   continue
   
        do 31 j = 1, my
            u(1,j) = 0.d0			! ��
            u(0,j) = u(1,j)
            v(0,j) = -v(1,j)
            p(0,j) = p(1,j)
            u(mx+1,j) = 0.d0		! �E
            u(mx+2,j) = u(mx+1,j)
            v(mx+1,j) = -v(mx,j)
            p(mx+1,j) = p(mx,j)
            p(mx+2,j) = p(mx+1,j)
   31   continue
   
!       u�̒�`�_��v4�Cv�̒�`�_�ɂ�����u4�iu��v���ꎞ���i�񂾂̂ŁCu4�Cv4���V��������j
        do 32 i = 1, mx
            do 32 j = 2, my
                u4(i,j) = (u(i,j) + u(i+1,j) + u(i,j-1) + u(i+1,j-1)) * 0.25d0
   32   continue
   
        do 33 i = 2, mx
            do 33 j = 1, my
                v4(i,j) = (v(i,j) + v(i-1,j) + v(i,j+1) + v(i-1,j+1)) * 0.25d0
   33   continue
   
        do 34 i = 1, mx
            u4(i,my+1) = u4(i,my)
   34   continue
   
        do 35 j = 1, my
            v4(mx+1,j) = -v4(mx,j)
!            v4(mx+1,j) = v4(mx,j)
   35   continue
   
!       ���͂̌v�Z
        do 40 i = 1, mx
            do 40 j = 1, my
                pp1(i,j) = ((u(i+1,j) - u(i,j))/dh + (v(i,j+1) - v(i,j))/dh)/dt

                pp2(i,j) = - ( &
                             ( &
						     ( u(i+1,j) + dabs( u(i+1,j)))/2.d0 * (u(i+1,j)   - u(i,j))/dh &
                          +  ( u(i+1,j) - dabs( u(i+1,j)))/2.d0 * (u(i+2,j)   - u(i+1,j))/dh &
            	    	  +  (v4(i+1,j) + dabs(v4(i+1,j)))/2.d0 * (u(i+1,j)   - u(i+1,j-1))/dh &
                          +  (v4(i+1,j) - dabs(v4(i+1,j)))/2.d0 * (u(i+1,j+1) - u(i+1,j))/dh &
                              ) &
                           -  ( &
			    			 ( u(i,j) + dabs( u(i,j)))/2.d0 * (u(i,j)   - u(i-1,j))/dh &
                          +  ( u(i,j) - dabs( u(i,j)))/2.d0 * (u(i+1,j) - u(i,j))/dh &
                	 	  +  (v4(i,j) + dabs(v4(i,j)))/2.d0 * (u(i,j)   - u(i,j-1))/dh &
                          +  (v4(i,j) - dabs(v4(i,j)))/2.d0 * (u(i,j+1) - u(i,j))/dh &
                              ) &
                              ) / dh

                pp3(i,j) = - ( &
                             ( &
						     (u4(i,j+1) + dabs(u4(i,j+1)))/2.d0 * (v(i,j+1)   - v(i-1,j+1))/dh &
                          +  (u4(i,j+1) - dabs(u4(i,j+1)))/2.d0 * (v(i+1,j+1) - v(i,j+1))/dh &
						  +  ( v(i,j+1) + dabs( v(i,j+1)))/2.d0 * (v(i,j+1)   - v(i,j))/dh &
                          +  ( v(i,j+1) - dabs( v(i,j+1)))/2.d0 * (v(i,j+2)   - v(i,j+1))/dh &
                             ) &
                           - ( &
			    			 (u4(i,j) + dabs(u4(i,j)))/2.d0 * (v(i,j)   - v(i-1,j))/dh &
                          +  (u4(i,j) - dabs(u4(i,j)))/2.d0 * (v(i+1,j) - v(i,j))/dh &
						  +  ( v(i,j) + dabs( v(i,j)))/2.d0 * (v(i,j)   - v(i,j-1))/dh &
                          +  ( v(i,j) - dabs( v(i,j)))/2.d0 * (v(i,j+1) - v(i,j))/dh &
                             ) &
                             ) / dh
                              
                phi(i,j) = rho * (pp1(i,j) + pp2(i,j) + pp3(i,j))
   40   continue
   
        knum = 0
  333   continue   
        do 41 i = 1, mx
            do 41 j = 1, my
!                pnew(i,j) = (1.d0 - alpha) * p(i,j) &
!                          + alpha * ((p(i+1,j) + p(i-1,j) + p(i,j+1) + p(i,j-1) &
!                          - 4.d0 * p(i,j)) / (dh**2) - phi(i,j))
                pnew(i,j) = ((p(i+1,j) + p(i-1,j) + p(i,j+1) + p(i,j-1)) &
                          - phi(i,j) * dh**2) / 4.d0
   41   continue

!------ ���͂̓��덷�̌v�Z -----------------------------------------------
        psum = 0.d0
        do 42 i = 1, mx
            do 42 j = 1, my
                psum = psum + (pnew(i,j) - p(i,j))**2  ! ���͂̓��덷�̘a���v�Z
   42   continue
        ! psum = psum / (dble(mx)*dble(my))

        if(psum > eps) then	! ���덷�̘a��eps���傫��������ȉ����J��Ԃ�
            do 43 i = 1, mx
                do 43 j = 1, my
                    p(i,j) = pnew(i,j)	! �V����pnew��p�ɏ㏑������
   43       continue

            knum = knum + 1	! �����v�Z�̉񐔂��J�E���g
            
            if(knum < kend) go to 333  ! �����v�Z�̃J�E���g������l�𒴂��Ȃ�������J��Ԃ�
 
        end if

!        write(*,*) 'knum= ',knum,' psum= ',psum,' eps= ',eps  ! �����v�Z�̌��ʂ�\��

!       ���ʂ̏o��

        if(mod(kt,1000) == 0) then

        write(*,*)
        write(*,*) 'kt = ',kt
        write(*,*) 'knum= ',knum,' psum= ',psum,' eps= ',eps  ! �����v�Z�̌��ʂ�\��
        
        read(10,'(a7)') file_name  ! �J�^���O�t�@�C����t�@�C������ǂ�
        open(20, file='./unew/unew_'//file_name,recl=200000)  ! ���O��unew�t�H���_������Ă���
        open(21, file='./vnew/vnew_'//file_name,recl=200000)  ! ���O��vnew�t�H���_������Ă���
        open(22, file='./pnew/pnew_'//file_name,recl=200000)  ! ���O��pnew�t�H���_������Ă���

        umax = -10.0
        umin = +10.0
        vmax = -10.0
        vmin = +10.0
        
        do 50 i = 1, mx
            do 50 j = 1, my
                uout(i,j) = (u(i+1,j) + u(i,j))/2.d0  ! �o�͗p�Ɋi�q���S�ł�u���v�Z
                vout(i,j) = (v(i,j+1) + v(i,j))/2.d0  ! �o�͗p�Ɋi�q���S�ł�v���v�Z

                if(uout(i,j) > umax) umax = uout(i,j)
                if(uout(i,j) < umin) umin = uout(i,j)
                if(vout(i,j) > vmax) vmax = vout(i,j)
                if(vout(i,j) < vmin) vmin = vout(i,j)        
   50   continue

        write(*,200) 'umax= ',umax,' umin= ',umin
        write(*,200) 'vmax= ',vmax,' vmin= ',vmin
  200   format(2(a6,f8.4))
   
        do 51 j = my, 1, -1
            write(20,100) (uout(i,j),i=1,mx)
            write(21,100) (vout(i,j),i=1,mx)
            write(22,100) (pnew(i,j),i=1,mx)
   51   continue
  100   format(f10.6,63(',',f10.6))

        close(20)	! �����ԍ��ōĂъJ���̂ŁC�K������
        close(21)	! �����ԍ��ōĂъJ���̂ŁC�K������
        close(22)	! �����ԍ��ōĂъJ���̂ŁC�K������

		end if

!       ���ʂ̓���
        do 60 i = 1, mx
            do 60 j = 1, my
      	        p(i,j) = pnew(i,j)  ! unew��vnew�͊��ɓ���ւ���Ă���̂ŁC���̂܂܎g��
   60   continue
   
! ----- ���Ԕ��W�̌v�Z�i�I���j------------------------------------
  500   continue

        stop
        end