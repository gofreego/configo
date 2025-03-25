import 'package:configo_ui/routes/route_constants.dart';
import 'package:configo_ui/widgets/main_layout.dart';
import '../screens/config/config_list.dart';
import 'package:go_router/go_router.dart';

final GoRouter appRouter = GoRouter(
  initialLocation: "/",
  routes: [
    ShellRoute(
      builder: (context, state, child) {
        return MainLayout(child: child);
      },
      routes: [
        GoRoute(
          name: RouteName.home,
          path: "/",
          builder: (context, state) {
            return ListConfigScreen();
          },
        ),
      ],
    ),
  ],
);
