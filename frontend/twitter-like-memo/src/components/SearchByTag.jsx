import React, { useState, useRef } from "react";
import { searchMemosByTags } from '../api';

export const SearchMemoByTags = ({ onSuccess }) => {
    const initialState = {
        tags: '',
    };

    const [values, setValues] = useState(initialState);
    const [searchResults, setSearchResults] = useState([]);
    const [hasSearched, setHasSearched] = useState(false);
    const tagsRef = useRef(null);

    const onValueChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        });
    };

    const handleSubmit = async (event) => {
        event.preventDefault();

        // タグが空の場合は検索を実行しない
        if (!values.tags.trim()) {
            alert('タグを入力してください');
            return;
        }

        setHasSearched(true);

        try {
            const results = await searchMemosByTags(values);
            setSearchResults(results);
            if (onSuccess) {
                onSuccess();
            }
        } catch (error) {
            console.error('GET error:', error);
            alert(error.message);
        }
    };

    const handleClear = () => {
        setValues(initialState);
        setSearchResults([]);
        setHasSearched(false);
        if (tagsRef.current) {
            tagsRef.current.value = '';
        }
    };

    return (
        <div>
            <div className="search-input-container">
                <input
                    ref={tagsRef}
                    className='input-tag'
                    type='text'
                    name="tags"
                    placeholder='タグ検索'
                    onChange={onValueChange}
                    value={values.tags}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            handleSubmit(e);
                        }
                    }}
                />
                <button
                    type="button"
                    onClick={handleSubmit}
                    className="search-button"
                >
                    検索
                </button>
                {hasSearched && (
                    <button
                        type="button"
                        onClick={handleClear}
                        className="clear-button"
                    >
                        X 結果をクリア
                    </button>
                )}
            </div>

            <div className="search-results">
            <p>{searchResults.length} 件見つかりました</p>
                {searchResults.length > 0 ? (
                    searchResults.map((result, index) => (
                        <div key={index} className="timeline">
                            
                            <div className='memo-card'>
                                <div className='memo-content'>
                                    <div className='memo-title'>{result.title}</div>
                                    <div className='memo-body'>{result.body}</div>
                                    {result.tags && (
                                        <div className="memo-tags">
                                            {result.tags.split(',').map((tag, tagIndex) => (
                                                <span key={tagIndex} className="tag">#{tag.trim()} </span>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    hasSearched && <p>検索結果がありません</p>
                )}
            </div>
        </div>
    );
};