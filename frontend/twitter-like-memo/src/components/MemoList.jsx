import { useEffect, useState } from "react";
import { fetchMemos, Memo } from "../api";

export const MemoList = (props) => {
    const [memos, setMemos] = useState([]);
    const reload = props.reload;

    useEffect(() => {
        fetchMemos()
            .then((data) => {
                console.debug('GET success:', data);
                setMemos(data.memos);
            })
            .catch((error) => {
                console.error('GET error:', error);
            });
    }, [reload]);

    return (
        <>
            {memos.map((memo) => (
                <div className='timeline' key={memo.id}>
                    <div className='memo-card'>
                        <div className='memo-title'>{memo.title}</div>
                        <div className='memo-body'>{memo.body}</div>
                    </div>
                    <hr />
                </div>
            ))}
        </>
    );
};