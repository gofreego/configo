import 'package:flutter/material.dart';

class ListConfigScreen extends StatefulWidget {
  const ListConfigScreen({super.key});

  @override
  State<ListConfigScreen> createState() => _ListConfigScreenState();
}

class _ListConfigScreenState extends State<ListConfigScreen> {
  @override
  Widget build(BuildContext context) {
    return const Center(
      child: AboutDialog(
        applicationIcon: FlutterLogo(),
        applicationName: 'Configo List',
        applicationVersion: '1.0.0',
        children: [
          Text('Configo is a simple configuration app.'),
        ],
      ),
    );
  }
}