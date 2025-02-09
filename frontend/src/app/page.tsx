import ContainerTable from "./components/ContainerTable";
import styles from "./page.module.css";

export default function Home() {
  return (
    <div className={styles.pageWrapper}>
      {}
      <header className={styles.header}>
        <h1 className={styles.headerTitle}>Docker Monitoring App</h1>
      </header>

      {}
      <section className={styles.tableSection}>
        <ContainerTable />
      </section>

      {}
      <footer className={styles.footer}>
        made by <strong>@k6zma</strong>
      </footer>
    </div>
  );
}
