import 'package:flutter/material.dart';

class MainLayout extends StatelessWidget {
  final Widget child;
  const MainLayout({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
       appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Image.asset("assets/images/logo-name.png", height: 40,),
      ),
      body: child,
    );
  }
}