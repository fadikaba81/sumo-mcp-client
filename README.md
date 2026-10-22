# sumo-mcp


+------------------+
| Go Application   |
| (MCP Client)     |
+--------+---------+
         |
         | MCP Protocol
         |
+--------v---------+
| Sumo Logic MCP   |
| Server           |
+--------+---------+
         |
         | REST API
         |
+--------v---------+
| Sumo Logic       |
| Search / Logs    |
| Metrics / Traces |
+------------------+