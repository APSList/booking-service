package booking

import "go.uber.org/fx"

// ======== EXPORTS ========

// Module exports services present
var Context = fx.Options(
	fx.Provide(GetReservationController),
	fx.Provide(
		fx.Annotate(
			GetReservationService,
			fx.As(new(Service)),
		),
	),
	fx.Provide(GetReservationRepository),
	fx.Provide(SetReservationRoutes),
)
