import type {NextPage} from "next";
import Content from "../posts/index.mdx";

const Home: NextPage = () => {
    return (
        <div>
            <div className="container m-auto max-w-screen-md p-4 pt-8 space-y-10">
                <article className="markdown-body leading-loose">
                    <Content/>
                </article>
            </div>
        </div>
    );
};

export default Home;
