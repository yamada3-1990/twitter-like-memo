import React, { useState, useRef } from "react";
import { addMemos } from './../api';

export const Adding = ({ onSuccess }) => {
    const initialState = {
        title: '',
        body: '',
        tags: '',
    };

    const [values, setValues] = useState(initialState);
    const titleRef = useRef(null);
    const bodyRef = useRef(null);
    const tagsRef = useRef(null);

    const onValueChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        });
    };

    const clearForm = () => {
        setValues(initialState);
        if (titleRef.current) titleRef.current.value = '';
        if (bodyRef.current) bodyRef.current.value = '';
        if (tagsRef.current) tagsRef.current.value = '';
    };

    const onSubmit = async (event) => {
        event.preventDefault();

        // バリデーション
        const REQUIRED_FIELDS = ['body'];
        const missingFields = Object.entries(values)
            .filter(([key, value]) => REQUIRED_FIELDS.includes(key) && !value)
            .map(([key]) => key);

        if (missingFields.length) {
            alert(`必須項目が入力されていません: ${missingFields.join(', ')}`);
            return;
        }

        console.log('Submitting values:', values);

        try {
            await addMemos(values);
            alert('メモを追加しました');
            clearForm();
            if (onSuccess) {
                onSuccess();
            }
        } catch (error) {
            console.error('POST error:', error);
            alert(error.message);
        }
    };

    return (
        <form onSubmit={onSubmit}>
            <div className='add-memo'>
                <input 
                    ref={titleRef}
                    className='input-title' 
                    type='text' 
                    name="title" 
                    placeholder='タイトル' 
                    onChange={onValueChange} 
                    value={values.title}
                />
                <input 
                    ref={bodyRef}
                    className='input-body' 
                    type='text' 
                    name="body" 
                    placeholder='本文' 
                    onChange={onValueChange}
                    value={values.body}
                />
                <input 
                    ref={tagsRef}
                    className='input-tag' 
                    type='text' 
                    name="tags" 
                    placeholder='タグ' 
                    onChange={onValueChange}
                    value={values.tags}
                />
                <button className='post-button' type="submit">メモを追加</button>
            </div>
        </form>
    );
};