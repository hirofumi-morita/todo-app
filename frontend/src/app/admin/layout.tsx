import type { Metadata } from "next";
import styles from "./adminLayout.module.css";

export const metadata: Metadata = {
  title: "管理画面 | TODO App",
  description: "管理画面のレイアウト",
};

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section className={styles.layout}>
      <div className={styles.content}>{children}</div>
    </section>
  );
}
