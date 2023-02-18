import styles from '../styles/Input.module.css'

export default function Input({ text, type }) {
    return (
        <div className={styles.inputWithLabel}>
            {text && <div className={styles.inputLabel}>{text}</div>}
            <input type={type} className={styles.inputBox} />
        </div>
    )
}