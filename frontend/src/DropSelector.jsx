import { useCallback, useEffect, useRef, useState } from 'react';
import { dropTargetForElements, monitorForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';

const DropSelector = ({ onDrop, itemType, dropText = 'Drag item over me', dropActiveText = 'Drop item here!' }) => {
    const ref = useRef(null);
    const [isDraggedOver, setIsDraggedOver] = useState(false);

    useEffect(() => {
        const el = ref.current;
        if (!el) {
            console.error("Draggable element not found.");
            return;
        }

        const dropTarget = dropTargetForElements({
            element: el,
            onDragEnter: () => setIsDraggedOver(true),
            onDragLeave: () => setIsDraggedOver(false),
            onDrop: (event) => {
                setIsDraggedOver(false);
                if (onDrop) {
                    onDrop(event);
                }
            }
        });

        return dropTarget;
    }, [onDrop]);

    return (
        <div
            className={`w-24 h-12 flex justify-center items-center border border-dashed rounded-md transition-all duration-300 ${isDraggedOver ? 'bg-blue-100' : 'bg-white'}`}
            ref={ref}
        >
            {isDraggedOver ? dropActiveText : dropText}
        </div>
    );
};

export default DropSelector;