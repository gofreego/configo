import 'package:flutter/material.dart';
import 'package:web/models/config/metadata.dart';

class ServiceInfoWidget extends StatelessWidget {
  
  final ServiceInfo serviceInfo;
  
  const ServiceInfoWidget({super.key, required this.serviceInfo});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          serviceInfo.name,
          style: const TextStyle(fontSize: 26, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 10),
        Text(serviceInfo.description, style: const TextStyle(fontSize: 16, color: Color.fromARGB(255, 63, 63, 63),fontWeight: FontWeight.w500),),
      ],
    );
  }
}
