import 'package:flutter/material.dart';

class ServiceInfoWidget extends StatelessWidget {
  final String serviceName = "Test Service";
  final String serviceDescription =
      "A powerful and scalable testing solution for APIs, applications, and security assessments. It provides end-to-end testing capabilities to ensure the reliability, security, and performance of your applications.";
  const ServiceInfoWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          serviceName,
          style: const TextStyle(fontSize: 26, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 10),
        Text(serviceDescription, style: const TextStyle(fontSize: 16, color: Color.fromARGB(255, 63, 63, 63),fontWeight: FontWeight.w500),),
      ],
    );
  }
}
