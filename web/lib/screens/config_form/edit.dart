import 'package:flutter/material.dart';

class EditConfigScreen extends StatefulWidget {
  final String id;
  const EditConfigScreen({super.key, required this.id});

  @override
  State<EditConfigScreen> createState() => _EditConfigScreenState();
}

class _EditConfigScreenState extends State<EditConfigScreen> {
  @override
  Widget build(BuildContext context) {
    return const Center(
      child: AboutDialog(
        applicationIcon: FlutterLogo(),
        applicationName: 'Configo Edit',
        applicationVersion: '1.0.0',
        children: [
          Text('Configo is a simple configuration app.'),
        ],
      ),
    );
  }
}