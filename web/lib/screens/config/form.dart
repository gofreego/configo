import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:web/models/config/object.dart';
import 'package:web/models/config/type.dart';
import 'package:flutter_code_editor/flutter_code_editor.dart';
import 'package:highlight/languages/json.dart';
import 'package:web/widgets/json_editor.dart';

class NewWidget extends StatelessWidget {
  const NewWidget({super.key});

  @override
  Widget build(BuildContext context) {
    List<ConfigObject> dummyConfigs = [
      ConfigObject(
        name: "Username",
        type: ConfigType.string,
        description: "Enter your username",
        required: true,
      ),
      ConfigObject(
        name: "About",
        type: ConfigType.bigText,
        description: "About yourself",
        required: true,
      ),
      ConfigObject(
        name: "Age",
        type: ConfigType.number,
        description: "Enter your age",
        required: true,
      ),
      ConfigObject(
        name: "Enable Notifications",
        type: ConfigType.boolean,
        description: "Turn on to receive notifications",
        required: false,
      ),
      ConfigObject(
        name: "Theme",
        type: ConfigType.choice,
        description: "Choose a theme",
        required: true,
        choices: ["Light", "Dark", "System"],
      ),
      // ConfigObject(
      //   name: "Settings",
      //   type: ConfigType.parent,
      //   description: "User preferences",
      //   required: false,
      //   children: [
      //     ConfigObject(
      //       name: "Font Size",
      //       type: ConfigType.number,
      //       description: "Adjust font size",
      //       required: false,
      //     ),
      //     ConfigObject(
      //       name: "Language",
      //       type: ConfigType.choice,
      //       description: "Select language",
      //       required: true,
      //       choices: ["English", "Spanish", "French"],
      //     ),
      //   ],
      // ),
      ConfigObject(
        name: "Custom JSON",
        type: ConfigType.json,
        description: "Enter custom configuration in JSON format",
        required: false,
      ),
    ];

    return ConfigForm(configs: dummyConfigs);
  }
}

class ConfigForm extends StatefulWidget {
  final List<ConfigObject> configs;

  const ConfigForm({super.key, required this.configs});

  @override
  _ConfigFormState createState() => _ConfigFormState();
}

class _ConfigFormState extends State<ConfigForm> {
  Map<String, dynamic> formValues = {};
  Map<String, bool> expandedStates = {};

  @override
  Widget build(BuildContext context) {
    List<ConfigObject> sortedConfigs = _sortConfigs(widget.configs);
    return LayoutBuilder(
      builder: (context, constraints) {
        return SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Wrap(
              spacing: 16.0,
              runSpacing: 16.0,
              children:
                  sortedConfigs
                      .map(
                        (config) =>
                            _buildConfigField(config, constraints.maxWidth),
                      )
                      .toList(),
            ),
          ),
        );
      },
    );
  }

  List<ConfigObject> _sortConfigs(List<ConfigObject> configs) {
    const typeOrder = {
      ConfigType.string: 1,
      ConfigType.choice: 2,
      ConfigType.number: 3,
      ConfigType.boolean: 4,
      ConfigType.bigText: 5,
      ConfigType.json: 6,
    };
    configs.sort(
      (a, b) => (typeOrder[a.type] ?? 99).compareTo(typeOrder[b.type] ?? 99),
    );
    return configs;
  }

  Widget _buildConfigField(ConfigObject config, double maxWidth) {
    if (config.children != null && config.children!.isNotEmpty) {
      return _buildExpandableConfig(config, maxWidth);
    }
    return _buildInputField(config, maxWidth);
  }

  Widget _buildExpandableConfig(ConfigObject config, double maxWidth) {
    return SizedBox(
      width: maxWidth,
      child: Card(
        margin: const EdgeInsets.symmetric(vertical: 8.0),
        child: Column(
          children: [
            ListTile(
              title: Text(
                config.name,
                style: const TextStyle(fontWeight: FontWeight.bold),
              ),
              subtitle: Text(config.description),
              trailing: Icon(
                expandedStates[config.name] == true
                    ? Icons.expand_less
                    : Icons.expand_more,
              ),
              onTap: () {
                setState(() {
                  expandedStates[config.name] =
                      !(expandedStates[config.name] ?? false);
                });
              },
            ),
            if (expandedStates[config.name] == true)
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16.0,
                  vertical: 8.0,
                ),
                child: Wrap(
                  spacing: 16.0,
                  runSpacing: 16.0,
                  children:
                      _sortConfigs(config.children!)
                          .map((child) => _buildConfigField(child, maxWidth))
                          .toList(),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildInputField(ConfigObject config, double maxWidth) {
    double fieldWidth = maxWidth > 400 ? 400 : maxWidth * 0.9;
    return SizedBox(
      width: fieldWidth,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            config.description,
            style: TextStyle(fontSize: 12, color: Colors.grey),
          ),
          const SizedBox(height: 4),
          _getFieldWidget(config),
        ],
      ),
    );
  }

  Widget _getFieldWidget(ConfigObject config) {
    switch (config.type) {
      case ConfigType.string:
        return TextFormField(
          initialValue: formValues[config.name] ?? '',
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          onChanged: (value) => formValues[config.name] = value,
        );
      case ConfigType.bigText:
        return TextFormField(
          initialValue: formValues[config.name] ?? '',
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          minLines: 3,
          maxLines: 5,
          onChanged: (value) => formValues[config.name] = value,
        );
      case ConfigType.number:
        return TextFormField(
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
          initialValue: formValues[config.name]?.toString() ?? '',
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          inputFormatters: [
            FilteringTextInputFormatter.allow(
              RegExp(r'^\d*\.?\d*$'),
            ), // Only allows numbers and a single decimal point
          ],
          onChanged: (value) {
            formValues[config.name] = num.tryParse(value) ?? 0;
          },
        );

      case ConfigType.boolean:
        return SwitchListTile(
          title: Text(config.name),
          value: formValues[config.name] ?? false,
          onChanged: (value) => setState(() => formValues[config.name] = value),
        );
      case ConfigType.choice:
        return DropdownButtonFormField<String>(
          value: formValues[config.name],
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          items:
              config.choices!
                  .map(
                    (choice) =>
                        DropdownMenuItem(value: choice, child: Text(choice)),
                  )
                  .toList(),
          onChanged: (value) => setState(() => formValues[config.name] = value),
        );
      case ConfigType.json:
        return JsonEditorWidget(name: config.name);
      default:
        return const SizedBox.shrink();
    }
  }
}
