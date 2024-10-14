import Clipboard from 'clipboard'
import { ElMessage } from 'element-plus'

export function copyClipboard(text, event) {
  const clipboard = new Clipboard(event.target, {
    text: () => text
  })
  clipboard.on('success', () => {
    ElMessage({
      type: "success",
      message: "复制成功"
    });
  })
  clipboard.onClick(event)
}
