import A from "./A";
import Head from "next/head";
import Login from "../images/login.svg"

const MainContainer = ({children, keywords}) => {
    return (
        <>
            <Head>
                <meta keywords={"ulbi tv, nextjs" + keywords}></meta>
                <title>Главная страница</title>
            </Head>
            <div className="navbar">
                <div><A href={'/'} text="LIGHTNESS"/></div>
                <div><span className="green-span">“Забота о здоровье - лучшая инвестиция в будущее!”</span></div>
                <a href="/auth"><Login /></a>
            </div>
            <div className="content">
                {children}
            </div>
            <style jsx>
                {`
                    .navbar {
                        display: flex;
                        height: 80px;
                        justify-content: space-around;
                        align-items: center;
                        font-family: Roboto, sans-serif;
                        box-shadow: 0px 3px 10px rgba(0, 0, 0, 0.05);
                    }
                   
                    .green-span {
                        color: #0A5F08;
                        font-size: 20px;
                    }

                    .content {
                        align-items: center;
                        display: flex;
                        justify-content: center;
                    }
                `}
            </style>
        </>
    );
};

export default MainContainer;
