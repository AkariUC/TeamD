import React from "react";
import Sidebar from "../components/sidebar";
import styles from '../css/home.module.css';

const HomePage = () => {
  return (
    <div className={styles.above}>
      <Sidebar/>
    </div>

  );
};

export default HomePage;
