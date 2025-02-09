"use client";

import { useQuery } from "@tanstack/react-query";
import { Table, Tag, Spin } from "antd";
import { fetchContainers, Container } from "../api/containers";

function compareIPs(ipA: string, ipB: string) {
  const octetsA = ipA.split(".").map(Number);
  const octetsB = ipB.split(".").map(Number);

  for (let i = 0; i < 4; i++) {
    const diff = octetsA[i] - octetsB[i];
    if (diff !== 0) {
      return diff;
    }
  }
  return 0;
}

const columns = [
  {
    title: "IP-адрес",
    dataIndex: "ip_address",
    key: "ip_address",
  },
  {
    title: "Имя контейнера",
    dataIndex: "name",
    key: "name",
  },
  {
    title: "Статус",
    dataIndex: "status",
    key: "status",
    render: (status: string) => {
      const color =
        status === "running"
          ? "green"
          : status === "paused"
          ? "blue"
          : status === "exited"
          ? "red"
          : "orange";
      return <Tag color={color}>{status.toUpperCase()}</Tag>;
    },
  },
  {
    title: "Пинг (мс)",
    dataIndex: "ping_time",
    key: "ping_time",
    render: (ping: number) => (ping === -1 ? "N/A" : `${ping.toFixed(2)} ms`),
  },
  {
    title: "Последний успешный пинг",
    dataIndex: "last_successful_ping",
    key: "last_successful_ping",
    render: (date: string) => (date ? new Date(date).toLocaleString() : "—"),
  },
];

export default function ContainerTable() {
  const { data, error, isLoading } = useQuery<Container[]>({
    queryKey: ["containers"],
    queryFn: fetchContainers,
    refetchInterval: 5000,
  });

  if (isLoading) {
    return (
      <div style={{ display: "flex", justifyContent: "center", marginTop: 20 }}>
        <Spin tip="Загрузка контейнеров..." size="large" />
      </div>
    );
  }

  if (error) {
    return <div>Ошибка загрузки данных</div>;
  }

  const sortedData = data ? [...data].sort((a, b) => compareIPs(a.ip_address, b.ip_address)) : [];

  return (
    <Table
      columns={columns}
      dataSource={sortedData}
      rowKey="container_id"
      pagination={false}
    />
  );
}
