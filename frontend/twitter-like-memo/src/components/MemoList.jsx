import { useEffect, useState } from "react";
import { fetchMemos, Memo } from "../api";
import { Delete } from "./Delete";

export const MemoList = (props) => {
    const [memos, setMemos] = useState([]);
    const reload = props.reload;
    const onDelete = props.onDelete;

    const fetchMemoList = async () => {
        try {
            const data = await fetchMemos();
            console.debug('GET success:', data);
            setMemos(data.memos || []);
        } catch (error) {
            console.error('GET error:', error);
            setMemos([]);
        }
    };

    useEffect(() => {
        fetchMemoList();
    }, [reload]);

    const handleDelete = () => {
        fetchMemoList(); // 削除後にメモリストを更新
        if (onDelete) {
            onDelete(); // 親コンポーネントに削除完了を通知
        }
    };

    return (
        <>
            {memos.map((memo) => (
                <div className='timeline' key={memo.id}>
                    <div className='memo-card'>
                        <div className='memo-header'>
                            <Delete memoId={memo.id} onDelete={handleDelete} />
                        </div>
                        <div className='memo-content'>
                            <div className='memo-title'>{memo.title}</div>
                            <div className='memo-body'>{memo.body}</div>
                            <div className='memo-tags'>
                                {memo.tags && memo.tags.split(',').map((tag, index) => (
                                    <span key={index} className='tag'>#{tag.trim()} </span>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </>
    );
};