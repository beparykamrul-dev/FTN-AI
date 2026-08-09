import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() => runApp(const FTNEnterpriseApp());

class FTNEnterpriseApp extends StatelessWidget {
  const FTNEnterpriseApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'FTN Enterprise Control',
        debugShowCheckedModeBanner: false,
        theme: ThemeData.dark(useMaterial3: true),
        home: const FTNHomePage(),
      );
}

class FTNHomePage extends StatefulWidget {
  const FTNHomePage({super.key});
  @override
  State<FTNHomePage> createState() => _FTNHomePageState();
}

class _FTNHomePageState extends State<FTNHomePage> {
  String status = 'Disconnected';
  List<dynamic> tools = const [];
  final endpoint = const String.fromEnvironment(
    'FTN_API_BASE_URL',
    defaultValue: 'http://127.0.0.1:5000',
  );

  Future<void> loadStatus() async {
    setState(() => status = 'Connecting...');
    try {
      final r = await http.get(
        Uri.parse('$endpoint/api/v1/ftn/android/status'),
        headers: const {'Accept': 'application/json'},
      );
      if (r.statusCode != 200) throw Exception('HTTP ${r.statusCode}');
      final data = jsonDecode(r.body) as Map<String, dynamic>;
      setState(() {
        status = '${data['service'] ?? 'FTN'} - ${data['status'] ?? 'unknown'}';
        tools = (data['tools'] as List?) ?? const [];
      });
    } catch (e) {
      setState(() => status = 'Connection failed: $e');
    }
  }

  @override
  void initState() {
    super.initState();
    loadStatus();
  }

  @override
  Widget build(BuildContext context) => Scaffold(
        appBar: AppBar(
          title: const Text('FTN Enterprise Control'),
          actions: [IconButton(onPressed: loadStatus, icon: const Icon(Icons.refresh))],
        ),
        body: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Card(child: ListTile(
              leading: const Icon(Icons.cloud_done),
              title: const Text('Backend'),
              subtitle: Text(status),
            )),
            const SizedBox(height: 16),
            const Text('Integrated Modules', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Expanded(
              child: tools.isEmpty
                  ? const Center(child: Text('No modules reported by backend'))
                  : ListView.builder(
                      itemCount: tools.length,
                      itemBuilder: (_, i) => Card(
                        child: ListTile(
                          leading: const Icon(Icons.extension),
                          title: Text('${tools[i]}'),
                          subtitle: const Text('Reported by FTN backend'),
                        ),
                      ),
                    ),
            ),
          ]),
        ),
      );
}
