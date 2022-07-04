package controller

func (s *SimulationController) horVeloBoundaryCondition() {
	// 最左2列の流速は0
	leftmost := 1
	for y := 1; y <= s.HorVeloCV.MaxHeight; y++ {
		s.HorVeloCV.Set(leftmost, y, 0.)
		s.HorVeloCV.Set(leftmost+1, y, 0.)
	}

	// 最右2列の流速は0
	rightmost := s.HorVeloCV.MaxWidth
	for y := 1; y <= s.HorVeloCV.MaxHeight; y++ {
		s.HorVeloCV.Set(rightmost, y, 0.)
		s.HorVeloCV.Set(rightmost-1, y, 0.)
	}

	// 最下2段の流速は上段と同じ
	bottom := 1
	for x := 1; x <= s.HorVeloCV.MaxWidth; x++ {
		v := s.HorVeloCV.Get(x, bottom+2)
		s.HorVeloCV.Set(x, bottom, v)
		s.HorVeloCV.Set(x, bottom+1, v)
	}

	// 最上2段の流速は主流
	top := s.HorVeloCV.MaxHeight
	for x := 1; x <= s.HorVeloCV.MaxWidth; x++ {
		s.HorVeloCV.Set(x, top, s.MainFlow)
		s.HorVeloCV.Set(x, top-1, s.MainFlow)
	}
}

func (s *SimulationController) verVeloBoundaryCondition() {
	// 最左2列の流速は右列と同じ
	leftmost := 1
	for y := 1; y <= s.VerVeloCV.MaxHeight; y++ {
		v := s.VerVeloCV.Get(leftmost+2, y)
		s.VerVeloCV.Set(leftmost, y, v)
		s.VerVeloCV.Set(leftmost+1, y, v)
	}

	// 最右2列の流速は左列と同じ
	rightmost := s.VerVeloCV.MaxWidth
	for y := 1; y <= s.VerVeloCV.MaxHeight; y++ {
		v := s.VerVeloCV.Get(rightmost-2, y)
		s.VerVeloCV.Set(rightmost, y, v)
		s.VerVeloCV.Set(rightmost-1, y, 0.)
	}

	// 最下2段の流速は0
	bottom := s.VerVeloCV.MinHeight
	for x := 1; x <= s.VerVeloCV.MaxWidth; x++ {
		s.VerVeloCV.Set(x, bottom, 0.)
		s.VerVeloCV.Set(x, bottom+1, 0.)
	}

	// 最上2段の流速は0
	top := s.VerVeloCV.MaxHeight
	for x := 1; x <= s.VerVeloCV.MaxWidth; x++ {
		s.VerVeloCV.Set(x, top, 0.)
		s.VerVeloCV.Set(x, top-1, 0.)
	}
}

func (s *SimulationController) pressBoundaryCondition() {
	// 最左2列の圧力は右列と同じ
	leftmost := s.PressCV.MinWidth
	for y := 1; y <= s.PressCV.MaxHeight; y++ {
		v := s.PressCV.Get(leftmost+2, y)
		s.PressCV.Set(leftmost, y, v)
		s.PressCV.Set(leftmost+1, y, v)
	}

	// 最右2列の圧力は左列と同じ
	rightmost := s.PressCV.MaxHeight
	for y := 1; y <= s.PressCV.MaxHeight; y++ {
		v := s.PressCV.Get(rightmost-2, y)
		s.PressCV.Set(rightmost, y, v)
		s.PressCV.Set(rightmost-1, y, v)
	}

	// 最下2段の圧力は上段と同じ
	bottom := s.PressCV.MinHeight
	for x := 1; x <= s.PressCV.MaxWidth; x++ {
		v := s.PressCV.Get(bottom+2, x)
		s.PressCV.Set(x, bottom, v)
		s.PressCV.Set(x, bottom+1, v)
	}

	// 最上2段の圧力は下段と同じ
	top := s.PressCV.MaxHeight
	for x := 1; x <= s.PressCV.MaxWidth; x++ {
		v := s.PressCV.Get(top-2, x)
		s.PressCV.Set(x, top, v)
		s.PressCV.Set(x, top-1, v)
	}
}
