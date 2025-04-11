import React, { useState, useRef } from "react";
import { deleteMemo } from '../api';

export const Delete = ({ memoId, onDelete }) => {
    const handleDelete = async () => {
        try {
            await deleteMemo(memoId);
            onDelete(); // 親コンポーネントに削除完了を通知
        } catch (error) {
            console.error('Delete error:', error);
            alert('メモの削除に失敗しました');
        }
    };

    return (
        <button className="delete-button" onClick={handleDelete}></button>
    );
};