'use client';

import { useEffect, useState } from 'react';
import { Alert } from '../lib/alert';
import { AreaChart, Card, Grid, Tab, TabList, Table, TableBody, TableCell, TableHead, TableHeaderCell, TableRow, Text, Title } from "@tremor/react";
import * as d3 from 'd3-array';

export default function AlertTable() {
  const [alerts, setAlerts] = useState<Alert[]>([]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:12000/');

    console.log('Connecting to websocket');

    ws.onopen = () => {
      console.log('Connected to websocket');
    };

    ws.onmessage = (event) => {
      console.log('Received message', event.data);
      const alert = JSON.parse(event.data);
      setAlerts((alerts) => [alert, ...alerts]);
    };

    return () => {
      console.log('Closing websocket');
      ws.close();
    };
  }, []);

  return (
    <main className="bg-slate-50 p-6 sm:p-10">
      <Title>Dashboard</Title>
      <Text>Fraud detection alerts</Text>

      <Grid numColsLg={2} className="mt-6 gap-6">
        <Card>
          <AreaChart
            data={Array.from(
              d3.sort(
                d3.rollups(
                  alerts, 
                  v => v.length,
                  d => Math.floor(new Date(d.transaction.utc).getTime() / 10000) * 10000,
                ).map(
                  ([k, v]) => 
                  ({seconds: k, Count: v, Time: new Date(k).toLocaleTimeString()})
                ),
                d => d.seconds
              )
            )}
            index="Time"
            categories={["Count"]}
            />
        </Card>
        <Card>
          <Table>
            <TableHead>
              <TableRow>
                <TableHeaderCell>Reason</TableHeaderCell>
                <TableHeaderCell>Amount</TableHeaderCell>
                <TableHeaderCell>Limit</TableHeaderCell>
                <TableHeaderCell>Owner</TableHeaderCell>
                <TableHeaderCell>Exp Date</TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {alerts.map((item) => (
                <TableRow key={item.transaction.utc + item.reason}>
                  <TableCell><Text>{item.reason}</Text></TableCell>
                  <TableCell><Text>{item.transaction.amount}</Text></TableCell>
                  <TableCell><Text>{item.transaction.limit_left}</Text></TableCell>
                  <TableCell><Text>{item.transaction.owner.first_name + ' ' + item.transaction.owner.last_name}</Text></TableCell>
                  <TableCell><Text>{item.transaction.card.exp_month + '/' + item.transaction.card.exp_year}</Text></TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      </Grid>
    </main>
  );
}

