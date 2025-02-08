import { useEffect, useState } from 'react';
import { monitorForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
import DropSelector from './DropSelector';

const DropSelectorContainer = ({ externalSelectList, setExternalSelectList, Type, title = 'Items', dropText, dropActiveText, handleSubmit }) => {
    const [droppedItems, setDroppedItems] = useState([]);

    useEffect(() => {
        return monitorForElements({
            onDrop({ source }) {
                console.log('Dropped source:', source.data);
                handleSubmit();
                if (source.data.type === Type) {
                    setDroppedItems((prevItems) => {
                        const isAlreadyDropped = prevItems.some(
                            (item) => item === source.data.item
                        );
                        if (!isAlreadyDropped) {
                            return [...prevItems, source.data.item];
                        }
                        return prevItems;
                    });
                }
            },
        });
    }, [Type]);

    useEffect(() => {
        setExternalSelectList(droppedItems);
    }, [droppedItems, setExternalSelectList]);

    return (
        <div className="p-4 flex flex-col gap-4">
            <h2 className="text-xl font-semibold mb-2">{title}</h2>
            <div className="flex gap-2">
                {droppedItems.map((item, index) => (
                    <div key={index} className="p-2 border border-gray-300 rounded-md bg-gray-100">
                        {item}
                    </div>
                ))}
                <DropSelector 
                    onDrop={({ source }) => {
                        if (source.data.type === Type) {
                            setDroppedItems((prevItems) => {
                                const isAlreadyDropped = prevItems.some(
                                    (item) => item === source.data.item
                                );
                                if (!isAlreadyDropped) {
                                    return [...prevItems, source.data.item];
                                }
                                return prevItems;
                            });
                        }
                    }}
                    itemType={Type.charAt(0).toUpperCase() + Type.slice(1)} dropText={dropText} dropActiveText={dropActiveText}
                />
            </div>
        </div>
    );
};

export default DropSelectorContainer;