import Input from '../components/Input'
import Button from '../components/Button'
import MainContainer from "../components/MainContainer";
import styles from '../styles/Auth.module.css'

export default function Auth() {
    return (
        <MainContainer keywords={"authorization"}>

        <div className={styles.formAuth}>
            <h1 className={styles.h1}>ЗАРЕГИСТРИРОВАТЬСЯ</h1>

            <div className={styles.form}>
                <Input text="имя" />
                <div className={styles.blockInForm}>
                    <Input text="возраст" />
                    <Input text="вес" />
                    <Input text="рост" />
                </div>
                <Input text="телефон" type="tel" />
                <Input text="эл.почта" type="email" />
            </div>
            <Button text="зарегистрироваться" />
            </div>
            </MainContainer>
    )
}