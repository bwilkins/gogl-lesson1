package main

import (
  "log"
  "github.com/go-gl/gl"
  "github.com/go-gl/glfw"
  "io/ioutil"
)

type Screen struct {
  w int
  h int
}

func LoadShader(filename string, shader_type gl.GLenum) (shader gl.Shader) {
  var shader_bytes []byte
  var shader_err error
  shader = gl.CreateShader(shader_type)
  shader_bytes, shader_err = ioutil.ReadFile(filename)
  if shader_err != nil {
    log.Fatal("Could not load shader from file: ", filename)
  }

  shader.Source(string(shader_bytes))
  shader.Compile()

  if shader.Get(gl.COMPILE_STATUS) == gl.FALSE {
    log.Printf("Compile error in shader %s:\n", filename)
    log.Fatal(shader.GetInfoLog())
  }

  return
}

func LoadShaderProgram() (program gl.Program) {
  var vertex_shader, fragment_shader gl.Shader


  program = gl.CreateProgram()

  vertex_shader = LoadShader("vertex_shader.txt", gl.VERTEX_SHADER)
  fragment_shader = LoadShader("fragment_shader.txt", gl.FRAGMENT_SHADER)
  program.AttachShader(vertex_shader)
  program.AttachShader(fragment_shader)

  program.Link()

  program.DetachShader(vertex_shader)
  program.DetachShader(fragment_shader)

  return
}

func LoadTriangle(program gl.Program) (gVAO gl.VertexArray) {
  var gVBO gl.Buffer
  var vertexData []gl.GLfloat
  var attrib_loc gl.AttribLocation

  gVAO = gl.GenVertexArray()
  gVAO.Bind()

  gVBO = gl.GenBuffer()
  gVBO.Bind(gl.ARRAY_BUFFER)

  vertexData = []gl.GLfloat{
      //    x     y     z
          0.0,  0.8,  0.0,
         -0.8, -0.8,  0.0,
          0.8, -0.8,  0.0,
  }

  gl.BufferData(gl.ARRAY_BUFFER, len(vertexData), vertexData, gl.STATIC_DRAW)

  attrib_loc = program.GetAttribLocation("vert")
  attrib_loc.EnableArray()
  attrib_loc.AttribPointer(3, gl.FLOAT, false, 0, nil)

  gVBO.Unbind(gl.ARRAY_BUFFER)
  clearVA()

  return
}

func clearVA() {
  gl.VertexArray(0).Bind()
}

func Render(program gl.Program, gVAO gl.VertexArray) {
  gl.ClearColor(0, 0, 0, 1)
  gl.Clear(gl.COLOR_BUFFER_BIT)

  program.Use()
  gVAO.Bind()
  gl.DrawArrays(gl.TRIANGLES, 0, 3)

  glfw.SwapBuffers()

  clearVA()
  gl.ProgramUnuse()

}

func main() {
  glfw_init_err := glfw.Init()
  if glfw_init_err != nil {
    log.Fatal("GLFW: shit broke")
  }
  defer glfw.Terminate()

  screensize := Screen{800, 600}

  glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
  glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
  glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)

  glfw_openwindow_err := glfw.OpenWindow(screensize.w, screensize.h, 8, 8, 8, 8, 0, 0, glfw.Windowed)
  if glfw_openwindow_err != nil {
    log.Fatal("GLFW: could not open window")
  }

  defer glfw.CloseWindow()

  glew_init_err := gl.Init()
  if glew_init_err != 0 {
    log.Fatal("GLEW: could not init")
  }

  maj, min, _ := glfw.GLVersion()
  if maj != 3 && min != 2 {
    log.Fatal("GL: Couldn't get GL 3.2 core profile")
  }

  log.Printf("OpenGL Version: %s", gl.GetString(gl.VERSION))
  log.Printf("GLSL Version: %s", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
  log.Printf("Vendor: %s", gl.GetString(gl.VENDOR))
  log.Printf("Renderer: %s", gl.GetString(gl.RENDERER))

  program := LoadShaderProgram()

  gVAO := LoadTriangle(program)

  for glfw.WindowParam(glfw.Opened) == gl.TRUE {
    Render(program, gVAO)
  }

}