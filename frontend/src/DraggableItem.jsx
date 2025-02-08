import { useCallback, useEffect, useRef, useState } from 'react';
import { draggable } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';

const DraggableItem = ({ item, type }) => {
  const ref = useRef(null);
  const [dragging, setDragging] = useState(false);

  useEffect(() => {
      const el = ref.current;
      if (!el) {
          console.error("Draggable element not found.");
          return;
      }

      return draggable({
          element: el,
          onDragStart: () => {
              setDragging(true);
          },
          onDragEnd: () => {
              setDragging(false);
          },
          getInitialData: () => {
              return { item, type };
          },
      });
  }, [item, type]);

  return (
      <div
          ref={ref}
          draggable
          className={`${
              dragging ? 'bg-[#171718] cursor-grabbing' : 'bg-[#171718] cursor-grab'
          } text-slate-50 p-3 rounded-lg shadow-lg transition-transform duration-300 ease-in-out hover:shadow-xl`}
          style={{
              transform: dragging ? 'scale(1.05)' : 'scale(1)',
              boxShadow: dragging
                  ? '0px 8px 16px rgba(0, 0, 0, 0.3)'
                  : '0px 4px 8px rgba(0, 0, 0, 0.15)',
          }}
      >
          {item}
      </div>
  );
};

export default DraggableItem;