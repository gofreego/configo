import 'package:web/routes/route_constants.dart';
import 'package:web/widgets/main_layout.dart';
import '../screens/config_form/edit.dart';
import '../screens/config_list/config_list.dart';
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
        GoRoute(
          name: RouteName.editConfig,
          path: "/edit/:id",
          builder: (context, state) {
            final id = state.pathParameters["id"];
            return EditConfigScreen(id: id!);
          },
        ),
      ],
    ),
  ],
);
