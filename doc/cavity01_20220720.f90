!  2次元キャビティフローの数値計算
!  2022.07.20
!  圧力のポアソン方程式の中の差分を中心差分から風上差分へ

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
		
        open(10,file = 'fname.prn')  ! カタログファイル
        
!       計算条件
        mx = 64
        my = 64
        xL = 0.02d0             ! キャビティサイズ [2cm]
        ktend = 10000            ! 時間発展計算の回数
        dh = xL / dble(mx)      ! メッシュサイズ

        dt = 0.001d0			! [sec]
        Us = 0.02d0  			! [m/s] キャビティの外の主流の流速
        eps = 10.d0**(-6)       ! 圧力の収束計算の閾値
        kend = 100000	 	    ! 圧力の収束計算の限界回数
        dvs = 10.d0**(-6)       ! 水の動粘性係数 [m^2/s]
        rho = 1000.d0           ! 水の密度 [kg/m^3]
		alpha = 0.5d0
 
 !      初期条件
        do 10 i = 1, mx+1	! 流下方向の流れの計算点は，1～mx+1
            do 10 j = 1, my
                u(i,j) = 0.d0
   10   continue
     
        do 11 i = 1, mx
            do 11 j = 1, my+1	! 横断方向の流れの計算点は，1～my+1
                v(i,j) = 0.d0
   11   continue
     
        do 12 i = 1, mx
            do 12 j = 1, my		! 横断方向の圧力の計算点．主流の中は計算しない
                p(i,j) = 0.4d0	! Re = 400の時の初期値
   12   continue
   
! ----- 時間発展の計算（メイン）------------------------------------
        do 500 kt = 1, ktend
!        write(*,*)
!        write(*,*) 'kt = ',kt

!       境界条件
!			主流の中の圧力は計算しない．
!			主流の下側の鉛直速度も計算しない．

        do 20 i = 1, mx
            u(i,my+1) = Us			! 主流の中の値
            v(i,my+1) = 0.d0		! 主流へは流れ込まない
            v(i,my+2) = 0.d0		! 主流の上側の壁の値
            p(i,my+1) = p(i,my)		! 主流の中の圧力は計算しない
            p(i,my+2) = p(i,my+1)	! 主流のさらに上側の値を与える
            u(i,0) = -u(i,1)			! 壁の下側
            v(i,1) = 0.d0			! 下側の壁の値
            p(i,0) = p(i,1)			! 壁の下側
   20   continue

        do 21 j = 1, my
            u(1,j) = 0.d0			! 左
            v(0,j) = -v(1,j)
            p(0,j) = p(1,j)
            u(mx+1,j) = 0.d0		! 右
            v(mx+1,j) = -v(mx,j)
            p(mx+1,j) = p(mx,j)
   21   continue
   
!       uの定義点のv4，vの定義点におけるu4
        do 22 i = 1, mx
            do 22 j = 2, my
                u4(i,j) = (u(i,j) + u(i+1,j) + u(i,j-1) + u(i+1,j-1)) * 0.25d0
   22   continue
   
        do 23 i = 2, mx
            do 23 j = 1, my
                v4(i,j) = (v(i,j) + v(i-1,j) + v(i,j+1) + v(i-1,j+1)) * 0.25d0
   23   continue
           
!       NS方程式の計算(unew)
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

!       NS方程式の計算(vnew)
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
                u(i,j) = unew(i,j)   ! NS方程式によって新しくなったunewをuに上書きする → uの時刻が一つ進む
   26   continue
   
        do 27 i = 1, mx
            do 27 j = 2, my
                v(i,j) = vnew(i,j)   ! NS方程式によって新しくなったvnewをvに上書きする → vの時刻が一つ進む
   27   continue
   
!------------------------------------------------------------------- ここまででuとvが更新された

!       境界条件（キャビティの中が一時刻進んだので，境界も新しくする）
        do 30 i = 1, mx
            u(i,my+1) = Us			! 上
            u(i,my+2) = u(i,my+1)
            v(i,my+1) = 0.d0
            v(i,my+2) = v(i,my+1)
            v(i,my+3) = v(i,my+2)
            p(i,my+1) = p(i,my)
            p(i,my+2) = p(i,my+1)
            u(i,0) = -u(i,1)			! 下
            v(i,1) = 0.d0
            v(i,0) = v(i,1)
            p(i,0) = p(i,1)
   30   continue
   
        do 31 j = 1, my
            u(1,j) = 0.d0			! 左
            u(0,j) = u(1,j)
            v(0,j) = -v(1,j)
            p(0,j) = p(1,j)
            u(mx+1,j) = 0.d0		! 右
            u(mx+2,j) = u(mx+1,j)
            v(mx+1,j) = -v(mx,j)
            p(mx+1,j) = p(mx,j)
            p(mx+2,j) = p(mx+1,j)
   31   continue
   
!       uの定義点のv4，vの定義点におけるu4（uとvが一時刻進んだので，u4，v4も新しくする）
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
   
!       圧力の計算
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

!------ 圧力の二乗誤差の計算 -----------------------------------------------
        psum = 0.d0
        do 42 i = 1, mx
            do 42 j = 1, my
                psum = psum + (pnew(i,j) - p(i,j))**2  ! 圧力の二乗誤差の和を計算
   42   continue
        ! psum = psum / (dble(mx)*dble(my))

        if(psum > eps) then	! 二乗誤差の和がepsより大きかったら以下を繰り返す
            do 43 i = 1, mx
                do 43 j = 1, my
                    p(i,j) = pnew(i,j)	! 新しいpnewをpに上書きする
   43       continue

            knum = knum + 1	! 収束計算の回数をカウント
            
            if(knum < kend) go to 333  ! 収束計算のカウントが上限値を超えなかったら繰り返す
 
        end if

!        write(*,*) 'knum= ',knum,' psum= ',psum,' eps= ',eps  ! 収束計算の結果を表示

!       結果の出力

        if(mod(kt,1000) == 0) then

        write(*,*)
        write(*,*) 'kt = ',kt
        write(*,*) 'knum= ',knum,' psum= ',psum,' eps= ',eps  ! 収束計算の結果を表示
        
        read(10,'(a7)') file_name  ! カタログファイからファイル名を読む
        open(20, file='./unew/unew_'//file_name,recl=200000)  ! 事前にunewフォルダを作っておく
        open(21, file='./vnew/vnew_'//file_name,recl=200000)  ! 事前にvnewフォルダを作っておく
        open(22, file='./pnew/pnew_'//file_name,recl=200000)  ! 事前にpnewフォルダを作っておく

        umax = -10.0
        umin = +10.0
        vmax = -10.0
        vmin = +10.0
        
        do 50 i = 1, mx
            do 50 j = 1, my
                uout(i,j) = (u(i+1,j) + u(i,j))/2.d0  ! 出力用に格子中心でのuを計算
                vout(i,j) = (v(i,j+1) + v(i,j))/2.d0  ! 出力用に格子中心でのvを計算

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

        close(20)	! 同じ番号で再び開くので，必ず閉じる
        close(21)	! 同じ番号で再び開くので，必ず閉じる
        close(22)	! 同じ番号で再び開くので，必ず閉じる

		end if

!       結果の入替
        do 60 i = 1, mx
            do 60 j = 1, my
      	        p(i,j) = pnew(i,j)  ! unewとvnewは既に入れ替わっているので，そのまま使う
   60   continue
   
! ----- 時間発展の計算（終了）------------------------------------
  500   continue

        stop
        end