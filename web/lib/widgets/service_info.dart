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
          style: const TextStyle(fontSize: 30, fontWeight: FontWeight.bold, color: Color.fromARGB(255, 70, 67, 67)),
        ),
        const SizedBox(height: 10),
        Text(serviceInfo.description, style: const TextStyle(fontSize: 16, color: Color.fromARGB(255, 63, 63, 63),fontWeight: FontWeight.w500),),
      ],
    );
  }
}
